section .text
global pangram
pangram:
        xor edx, edx
.loop:
        movzx ecx, byte [rdi]
        cmp ecx,0                   ; if null character, jump to end
        jz .end
        or ecx, 32
        sub ecx, 97
        bts edx, ecx
        inc rdi
        jmp .loop
.end:
        xor eax, eax
        and edx, 0x03ffffff
        cmp edx, 0x03ffffff
        sete al
        ret