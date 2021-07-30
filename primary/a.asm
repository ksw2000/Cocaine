global main
extern printf

section .data
    main.year dq  2021
    main.$var1 db  "The primary Cocaine Compiler is created in %d", 0

section .text
main:
    push    rcx
    push    rdx
    sub     rsp, 32
    mov     rcx, main.$var1
    mov     rdx, [main.year]
    call    printf
    add     rsp, 32
    pop     rdx
    pop     rcx
    
    mov     qword [main.year], 2022
    push    rcx
    push    rdx
    sub     rsp, 32
    mov     rcx, main.$var1
    mov     rdx, [main.year]
    call    printf
    add     rsp, 32
    pop     rdx
    pop     rcx

; name -f win64 a.asm
; gcc a.o -o a.exe 