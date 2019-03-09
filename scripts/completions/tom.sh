# bash completion for tom                                  -*- shell-script -*-

__tom_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__tom_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__tom_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__tom_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__tom_handle_reply()
{
    __tom_debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __tom_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __tom_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
		if declare -F __tom_custom_func >/dev/null; then
			# try command name qualified custom func
			__tom_custom_func
		else
			# otherwise fall back to unqualified for compatibility
			declare -F __custom_func >/dev/null && __custom_func
		fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__tom_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__tom_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__tom_handle_flag()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __tom_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __tom_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __tom_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if __tom_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__tom_handle_noun()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __tom_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __tom_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__tom_handle_command()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_tom_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __tom_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__tom_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __tom_handle_reply
        return
    fi
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __tom_handle_flag
    elif __tom_contains_word "${words[c]}" "${commands[@]}"; then
        __tom_handle_command
    elif [[ $c -eq 0 ]]; then
        __tom_handle_command
    elif __tom_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __tom_handle_command
        else
            __tom_handle_noun
        fi
    else
        __tom_handle_noun
    fi
    __tom_handle_word
}

__gotime_projects_get()
{
    local -a projects
    readarray -t COMPREPLY < <(gotime projects 2>/dev/null | grep "$cur" | sed -e 's/ /\\ /g')
}

__gotime_get_projects()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __gotime_projects_get ""
	else
	    __gotime_projects_get ${nouns[${#nouns[@]} -1]}
    fi
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        gotime_start)
            __gotime_get_projects
            return
            ;;
        *)
            ;;
    esac
}

_tom_cancel()
{
    last_command="tom_cancel"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_config_set()
{
    last_command="tom_config_set"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("activity.stop_on_start")
    must_have_one_noun+=("backup.directory")
    must_have_one_noun+=("backup.max_to_keep")
    must_have_one_noun+=("data_dir")
    must_have_one_noun+=("projects.create_missing")
    noun_aliases=()
}

_tom_config()
{
    last_command="tom_config"

    command_aliases=()

    commands=()
    commands+=("set")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create_project()
{
    last_command="tom_create_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--parent=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--parent=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create_tag()
{
    last_command="tom_create_tag"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create()
{
    last_command="tom_create"

    command_aliases=()

    commands=()
    commands+=("project")
    commands+=("tag")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit_frame()
{
    last_command="tom_edit_frame"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--end=")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--end=")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--notes=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--notes=")
    flags+=("--project=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project=")
    flags+=("--start=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--start=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit_project()
{
    last_command="tom_edit_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--hourly-rate=")
    local_nonpersistent_flags+=("--hourly-rate=")
    flags+=("--name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name=")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--parent=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--parent=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit()
{
    last_command="tom_edit"

    command_aliases=()

    commands=()
    commands+=("frame")
    commands+=("project")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_frames_archive()
{
    last_command="tom_frames_archive"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--include-subprojects")
    local_nonpersistent_flags+=("--include-subprojects")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--project=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_frames()
{
    last_command="tom_frames"

    command_aliases=()

    commands=()
    commands+=("archive")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--delimiter=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter=")
    flags+=("--format=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format=")
    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--project=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project=")
    flags+=("--subprojects")
    flags+=("-s")
    local_nonpersistent_flags+=("--subprojects")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import_fanurio()
{
    last_command="tom_import_fanurio"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import_macTimeTracker()
{
    last_command="tom_import_macTimeTracker"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import_watson()
{
    last_command="tom_import_watson"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import()
{
    last_command="tom_import"

    command_aliases=()

    commands=()
    commands+=("fanurio")
    commands+=("macTimeTracker")
    commands+=("watson")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_invoice_sevdesk()
{
    last_command="tom_invoice_sevdesk"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--key=")
    two_word_flags+=("-k")
    local_nonpersistent_flags+=("--key=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    flags+=("--dry-run")
    flags+=("-d")
    flags+=("--month=")
    flags+=("--project=")
    two_word_flags+=("-p")
    flags+=("--round-frames=")
    flags+=("--round-frames-to=")

    must_have_one_flag=()
    must_have_one_flag+=("--key=")
    must_have_one_flag+=("-k")
    must_have_one_noun=()
    noun_aliases=()
}

_tom_invoice()
{
    last_command="tom_invoice"

    command_aliases=()

    commands=()
    commands+=("sevdesk")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dry-run")
    flags+=("-d")
    flags+=("--month=")
    flags+=("--project=")
    two_word_flags+=("-p")
    flags+=("--round-frames=")
    flags+=("--round-frames-to=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_flag+=("--project=")
    must_have_one_flag+=("-p")
    must_have_one_noun=()
    noun_aliases=()
}

_tom_projects()
{
    last_command="tom_projects"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter=")
    flags+=("--format=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format=")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--recent=")
    local_nonpersistent_flags+=("--recent=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove_all()
{
    last_command="tom_remove_all"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("all")
    must_have_one_noun+=("frames")
    must_have_one_noun+=("projects")
    must_have_one_noun+=("tags")
    noun_aliases=()
}

_tom_remove_frame()
{
    last_command="tom_remove_frame"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove_project()
{
    last_command="tom_remove_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove()
{
    last_command="tom_remove"

    command_aliases=()

    commands=()
    commands+=("all")
    commands+=("frame")
    commands+=("project")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_rename()
{
    last_command="tom_rename"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_report()
{
    last_command="tom_report"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--day=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--day=")
    flags+=("--decimal")
    local_nonpersistent_flags+=("--decimal")
    flags+=("--description=")
    local_nonpersistent_flags+=("--description=")
    flags+=("--from=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--from=")
    flags+=("--include-archived")
    local_nonpersistent_flags+=("--include-archived")
    flags+=("--json")
    local_nonpersistent_flags+=("--json")
    flags+=("--matrix-tables")
    local_nonpersistent_flags+=("--matrix-tables")
    flags+=("--month")
    flags+=("-m")
    local_nonpersistent_flags+=("--month")
    flags+=("--output-file=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output-file=")
    flags+=("--project=")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project=")
    flags+=("--round-frames=")
    local_nonpersistent_flags+=("--round-frames=")
    flags+=("--round-frames-to=")
    local_nonpersistent_flags+=("--round-frames-to=")
    flags+=("--round-totals=")
    local_nonpersistent_flags+=("--round-totals=")
    flags+=("--round-totals-to=")
    local_nonpersistent_flags+=("--round-totals-to=")
    flags+=("--save-config=")
    local_nonpersistent_flags+=("--save-config=")
    flags+=("--show-empty")
    local_nonpersistent_flags+=("--show-empty")
    flags+=("--show-sales")
    local_nonpersistent_flags+=("--show-sales")
    flags+=("--show-summary")
    local_nonpersistent_flags+=("--show-summary")
    flags+=("--show-tracked")
    local_nonpersistent_flags+=("--show-tracked")
    flags+=("--show-untracked")
    local_nonpersistent_flags+=("--show-untracked")
    flags+=("--split=")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--split=")
    flags+=("--subprojects")
    local_nonpersistent_flags+=("--subprojects")
    flags+=("--template=")
    local_nonpersistent_flags+=("--template=")
    flags+=("--template-file=")
    flags_with_completion+=("--template-file")
    flags_completion+=("__tom_handle_filename_extension_flag gohtml")
    local_nonpersistent_flags+=("--template-file=")
    flags+=("--title=")
    local_nonpersistent_flags+=("--title=")
    flags+=("--to=")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--to=")
    flags+=("--year=")
    two_word_flags+=("-y")
    local_nonpersistent_flags+=("--year=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_start()
{
    last_command="tom_start"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--create-missing")
    local_nonpersistent_flags+=("--create-missing")
    flags+=("--notes=")
    local_nonpersistent_flags+=("--notes=")
    flags+=("--stop-on-start")
    local_nonpersistent_flags+=("--stop-on-start")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_status_projects()
{
    last_command="tom_status_projects"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--delimiter=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter=")
    flags+=("--format=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format=")
    flags+=("--include-active")
    local_nonpersistent_flags+=("--include-active")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--show-empty")
    flags+=("-e")
    local_nonpersistent_flags+=("--show-empty")
    flags+=("--show-overall")
    local_nonpersistent_flags+=("--show-overall")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_status()
{
    last_command="tom_status"

    command_aliases=()

    commands=()
    commands+=("projects")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter=")
    flags+=("--format=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format=")
    flags+=("--name-delimiter=")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--verbose")
    flags+=("-v")
    local_nonpersistent_flags+=("--verbose")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_stop()
{
    last_command="tom_stop"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--all")
    flags+=("-a")
    local_nonpersistent_flags+=("--all")
    flags+=("--notes=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--notes=")
    flags+=("--past=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--past=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_tags()
{
    last_command="tom_tags"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter=")
    flags+=("--format=")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format=")
    flags+=("--output=")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output=")
    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_root_command()
{
    last_command="tom"

    command_aliases=()

    commands=()
    commands+=("cancel")
    commands+=("config")
    commands+=("create")
    commands+=("edit")
    commands+=("frames")
    commands+=("import")
    commands+=("invoice")
    commands+=("projects")
    commands+=("remove")
    commands+=("rename")
    commands+=("report")
    commands+=("start")
    commands+=("status")
    commands+=("stop")
    commands+=("tags")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    flags+=("--config=")
    two_word_flags+=("-c")
    flags+=("--data-dir=")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_tom()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __tom_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("tom")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __tom_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_tom tom
else
    complete -o default -o nospace -F __start_tom tom
fi

# ex: ts=4 sw=4 et filetype=sh
