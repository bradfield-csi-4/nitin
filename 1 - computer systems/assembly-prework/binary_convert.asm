section .text
global binary_convert
binary_convert:
    xor esi, esi
    jmp .test
.loop:
    inc rdi
    inc esi
.test:
    cmp byte [rdi],0
    jne .loop
    dec rdi
    xor rax,rax
    mov rdx,1
.loop2:
    cmp byte [rdi],49
    jne .false
    add rax,rdx
.false:
    shl rdx,1
    dec rdi
    dec esi
.test2:
    cmp esi,0
    jne .loop2
    ret