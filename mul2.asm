main:
    imm 5;
    mov E, A;
    imm mul2;
    call;
    mov out, F;
    imm end;
    jmp;
mul2:
    mov B, E;
    mov C, E;
    add;
    mov F, D;
    ret;
end:
    


