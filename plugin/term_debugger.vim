let g:term_loc_variable = '$LOC'
let g:term_debugger_bin = expand('<sfile>:p:h:h') . '/bin/term_debugger'

highlight TermDebuggerBreakpoint ctermfg=red guifg=red
highlight TermDebuggercursor ctermfg=green guifg=green
sign define TermDebuggerCursor text==> linehl=CursorLine texthl=TermDebuggerCursor
sign define TermDebuggerBreakpoint text=● texthl=TermDebuggerBreakpoint

function Tapi_TermDebuggerStepTo(bufnum, arglist)
	let l:filename = a:arglist[0]
	let l:line = a:arglist[1]

	if filereadable(l:filename)
		let l:key = t:cursor_filename . ':' . t:cursor_line
		if has_key(t:breakpoints, l:key)
			call s:placeBreakpoint(t:cursor_filename, t:cursor_line)
		else
			call sign_unplace('TermDebuggerCursor')
		endif
		2wincmd w
		" open the file stepped to
		exe 'edit +' . l:line . ' ' . l:filename
		let l:bufnr = winbufnr(2)
		z.
		" add cursor
		call sign_place(0, 'TermDebuggerCursor', 'TermDebuggerCursor', l:bufnr, {'lnum': l:line, 'priority': 90})
		let t:cursor_filename = l:filename
		let t:cursor_line = l:line
		1wincmd w
	else
		echom "file not found: " . l:filename . 'please manually cd to the correct directory'
	endif
endfunction

function Tapi_TermDebuggerSignBreakpoint(bufnum, arglist)
	" place breakpoint at l:filename l:line
	let l:brk_type = a:arglist[0]
	let l:filename = a:arglist[1]
	let l:line = a:arglist[2]

	let l:key = l:filename . ':' . l:line
	let t:breakpoints[l:key] = l:brk_type
	if l:filename != t:cursor_filename || l:line != t:cursor_line
		call term_debugger#placeBreakpoint(l:filename, l:line)
	endif
endfunction

command -nargs=+ -complete=file TermDebugger call term_debugger#open_terminal(g:term_debugger_bin . ' ' . <q-args>)
command -nargs=+ TermSendKeys call term_debugger#term_sendkeys(<q-args>)
