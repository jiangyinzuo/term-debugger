function term_debugger#term_exit_cb(job, exit_status)
	call sign_unplace('TermDebuggerCursor')
	call sign_unplace('TermDebuggerBreakpoint')
endfunction

function term_debugger#open_terminal(command)
	" filename:line -> brk_type
	let t:breakpoints = {}
	let t:cursor_filename = ""
	let t:cursor_line = 0
	let t:term_bufnr = term_start(a:command, {'exit_cb': function('term_debugger#term_exit_cb'), 'term_rows': 15})
endfunction

function term_debugger#term_sendkeys(command)
	if !exists('t:term_bufnr')
		echoerr "run :TermDebugger first"
		return
	endif
	let l:cmd = substitute(a:command, g:term_loc_variable, expand('%'), '')
	call term_sendkeys(t:term_bufnr, l:cmd . "\<CR>")
endfunction

function term_debugger#placeBreakpoint(filename, line)
	call bufadd(a:filename)
	let l:bnr = bufnr(a:filename)
	call sign_place(a:line, 'TermDebuggerBreakpoint', 'TermDebuggerBreakpoint',
				\ l:bnr, {'lnum': a:line, 'priority': 90})
endfunction

