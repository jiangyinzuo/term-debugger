#!/bin/bash
printf "\x1b]51;%s\x07" '["call","term_debugger#add_breakpoint",["hello.cpp",18]]'
printf "\x1b]51;%s\x07" '["call","Tapi_TestEscSeq",[]]'
#printf "\x1b]51;%s\x07" '["drop","README.md"]'
