section .text
global fib
fib:
    mov eax, edi
    cmp edi, 1
    jle .end
    push rbx
    push rbp
    lea rdi, [rdi - 1]
    mov ebx, edi
    call fib
    mov rbp, rax
    lea rdi, [rbx - 1]
    call fib
    add rax, rbp
    pop rbp
    pop rbx
.end:
	ret
