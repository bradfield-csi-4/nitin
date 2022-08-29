default rel

section .text
global volume
volume:
    mulss xmm0, xmm0
    mulss xmm0, xmm1
    mulss xmm0, [pi3]
 	ret

section .rodata
pi3: dd 1.0471975512