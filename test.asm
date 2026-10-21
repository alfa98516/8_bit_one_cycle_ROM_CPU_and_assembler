main:
    imm 4;
    mov E, A;
    imm loop;
    jmp;
loop:
    mov B, E;
    imm 1;
    mov C, A;
    sub;
    imm loop_end;
    jz;
    mov E, D;
    imm loop;
    jmp;
loop_end:
