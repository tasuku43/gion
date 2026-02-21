package cli

import (
	"fmt"
	"io"
	"os"
)

const SupportedShells = "bash, zsh"

func runCompletion(args []string) error {
	shell := "bash"
	if len(args) > 0 {
		shell = args[0]
	}
	switch shell {
	case "bash":
		printBashCompletion(os.Stdout)
		return nil
	case "zsh":
		printZshCompletion(os.Stdout)
		return nil
	default:
		return fmt.Errorf("unsupported shell: %s (supported: %s)", shell, SupportedShells)
	}
}

func printBashCompletion(w io.Writer) {
	fmt.Fprintln(w, `# gion bash completion
_gion_completion() {
  local cur prev words cword
  _init_completion || return

  local commands="init doctor repo manifest plan import apply version help completion"
  local manifest_subcmds="ls add rm gc validate preset"
  local manifest_aliases="man m"
  local preset_subcmds="ls add rm validate"
  local preset_aliases="pre p"
  local repo_subcmds="get ls rm"

  if [[ ${cword} -eq 1 ]]; then
    COMPREPLY=($(compgen -W "${commands} ${manifest_aliases}" -- "${cur}"))
    return
  fi

  case ${words[1]} in
    manifest|man|m)
      if [[ ${cword} -eq 2 ]]; then
        COMPREPLY=($(compgen -W "${manifest_subcmds} ${preset_aliases}" -- "${cur}"))
        return
      fi
      case ${words[2]} in
        preset|pre|p)
          if [[ ${cword} -eq 3 ]]; then
            COMPREPLY=($(compgen -W "${preset_subcmds}" -- "${cur}"))
            return
          fi
          case ${words[3]} in
            add)
              COMPREPLY=($(compgen -W "--repo --no-prompt" -- "${cur}"))
              return
            ;;
            rm)
              COMPREPLY=($(compgen -W "--no-prompt" -- "${cur}"))
              return
            ;;
          esac
        ;;
        add)
          COMPREPLY=($(compgen -W "--preset --review --issue --repo --branch --base --no-apply --no-prompt" -- "${cur}"))
          return
        ;;
        rm)
          COMPREPLY=($(compgen -W "--no-apply --no-prompt" -- "${cur}"))
          return
        ;;
        gc)
          COMPREPLY=($(compgen -W "--no-apply --no-fetch --no-prompt" -- "${cur}"))
          return
        ;;
      esac
    ;;
    repo)
      if [[ ${cword} -eq 2 ]]; then
        COMPREPLY=($(compgen -W "${repo_subcmds}" -- "${cur}"))
        return
      fi
      case ${words[2]} in
        rm)
          COMPREPLY=($(compgen -W "--no-prompt" -- "${cur}"))
          return
        ;;
      esac
    ;;
    doctor)
      COMPREPLY=($(compgen -W "--fix --self" -- "${cur}"))
      return
    ;;
    completion)
      if [[ ${cword} -eq 2 ]]; then
        COMPREPLY=($(compgen -W "bash zsh" -- "${cur}"))
        return
      fi
    ;;
  esac

  if [[ ${cur} == -* ]]; then
    COMPREPLY=($(compgen -W "--root --no-prompt --debug --help --version" -- "${cur}"))
  fi
}

complete -F _gion_completion gion`)
}

func printZshCompletion(w io.Writer) {
	fmt.Fprintln(w, `#compdef gion

_gion() {
  local -a commands
  commands=(
    'init:initialize root layout'
    'doctor:check workspace/repo health'
    'repo:repo commands'
    'manifest:manifest inventory commands'
    'plan:show manifest diff'
    'import:rebuild manifest from filesystem'
    'apply:apply manifest to filesystem'
    'version:print version'
    'help:show help'
    'completion:generate shell completion'
    'man:alias for manifest'
    'm:alias for manifest'
  )

  local -a repo_subcmds
  repo_subcmds=(
    'get:fetch or update bare repo store'
    'ls:list known bare repo stores'
    'rm:remove bare repo stores'
  )

  local -a manifest_subcmds
  manifest_subcmds=(
    'ls:list workspace inventory'
    'add:add workspace to manifest'
    'rm:remove workspace entries'
    'gc:garbage collect safe workspaces'
    'validate:validate manifest inventory'
    'preset:preset inventory commands'
    'pre:alias for preset'
    'p:alias for preset'
  )

  local -a preset_subcmds
  preset_subcmds=(
    'ls:list manifest presets'
    'add:add a preset entry'
    'rm:remove preset entries'
    'validate:validate presets'
  )

  local -a global_flags
  global_flags=(
    '--root[override root]:path'
    '--no-prompt[disable interactive prompt]'
    '--debug[write debug logs to file]'
    '--help[show help]'
    '--version[print version]'
  )

  _arguments -C "${global_flags[@]}" ':command:->command' '*::args:->args'

  case $state in
    command)
      _describe 'command' commands
    ;;
    args)
      case ${words[1]} in
        manifest|man|m)
          case ${words[2]} in
            preset|pre|p)
              case ${words[3]} in
                add)
                  _arguments '--repo[repo spec]' '--no-prompt[disable interactive prompt]'
                ;;
                rm)
                  _arguments '--no-prompt[disable interactive prompt]'
                ;;
                *)
                  _describe 'preset subcommand' preset_subcmds
                ;;
              esac
            ;;
            add)
              _arguments \
                '--preset[preset name]:name' \
                '--review[add review workspace from PR]:url' \
                '--issue[add issue workspace from issue]:url' \
                '--repo[add workspace from repo]:repo' \
                '--branch[override branch name]:name' \
                '--base[override base ref]:ref' \
                '--no-apply[update manifest only]' \
                '--no-prompt[disable interactive prompt]'
            ;;
            rm)
              _arguments '--no-apply[update manifest only]' '--no-prompt[disable interactive prompt]'
            ;;
            gc)
              _arguments '--no-apply[update manifest only]' '--no-fetch[disable git fetch]' '--no-prompt[disable interactive prompt]'
            ;;
            *)
              _describe 'manifest subcommand' manifest_subcmds
            ;;
          esac
        ;;
        repo)
          case ${words[2]} in
            get)
            ;;
            ls)
            ;;
            rm)
              _arguments '--no-prompt[disable interactive prompt]'
            ;;
            *)
              _describe 'repo subcommand' repo_subcmds
            ;;
          esac
        ;;
        doctor)
          _arguments '--fix[list issues and planned fixes]' '--self[run self-diagnostics]'
        ;;
        completion)
          _values 'shell' bash zsh
        ;;
      esac
    ;;
  esac
}

_gion`)
}
