let g:term_loc_variable = '$LOC'

sign define TermDebuggerCursor text==> linehl=CursorLine texthl=TermDebuggerCursor
sign define TermDebuggerBrkD text=○ texthl=TermDebuggerBrkD
sign define TermDebuggerBrk1 text=● texthl=TermDebuggerBrk1
sign define TermDebuggerBrk2 text=●² texthl=TermDebuggerBrk2
sign define TermDebuggerBrk3 text=●³ texthl=TermDebuggerBrk3
sign define TermDebuggerBrk4 text=●⁴ texthl=TermDebuggerBrk4
sign define TermDebuggerBrk5 text=●⁵ texthl=TermDebuggerBrk5
sign define TermDebuggerBrk6 text=●⁶ texthl=TermDebuggerBrk6
sign define TermDebuggerBrk7 text=●⁷ texthl=TermDebuggerBrk7
sign define TermDebuggerBrk8 text=●⁸ texthl=TermDebuggerBrk8
sign define TermDebuggerBrk9 text=●⁹ texthl=TermDebuggerBrk9
sign define TermDebuggerBrkN text=●ⁿ texthl=TermDebuggerBrkN

function Tapi_TermDebuggerStepTo(bufnum, arglist)
	let l:filename = a:arglist[1]
	let l:line = a:arglist[2]
	let l:cmd = "TermDebuggerStepTo " . l:filename . " " . l:line

endfunction

function Tapi_TermDebuggerSetBreakpoint(bufnum, arglist)
	let l:line = a:arglist[2]
	let l:filename = a:arglist[1]
	if a:arglist[0] == ' '
		call sign_unplace('TermDebuggerBrk', {'buffer': l:filename, 'id': l:line})
		echom "Removing breakpoint from " . l:filename . " at line " . l:line
	else	
		call sign_place(l:line, 'TermDebuggerBrk', 'TermDebuggerBrk' . a:arglist[0], l:filename, {'lnum': l:line, 'priority': 99}})
		echom "Adding " . a:arglist[0] . " breakpoint to " .  l:filename . " at line " . l:line
	endif
endfunction

function term_debugger#open_terminal(command)
	let t:term_filename = ''
	let t:term_line = 0
	let t:term_bufnr = term_start(a:command)
endfunction

function term_debugger#term_sendkeys(command)
	if !exists('t:term_bufnr')
		echoerr "run :TermDebugger first"
		return
	endif
	let l:cmd = substitute(a:command, g:term_loc_variable, expand('%'), '')
	call term_sendkeys(t:term_bufnr, l:cmd . "\<CR>")
endfunction

command -nargs=+ TermDebugger :call term_debugger#open_terminal(<q-args>)
command -nargs=+ TermSendKeys :call term_debugger#term_sendkeys(<q-args>)
