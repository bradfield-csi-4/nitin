section .text
global sum_to_n
sum_to_n:
    xor rax, rax
    xor rsi, rsi
.L1:
    add rax, rsi
    inc rsi
    cmp rsi, rdi
    jle .L1
    ret