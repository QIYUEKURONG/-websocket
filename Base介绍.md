# Base家族的介绍

## Base64 简介

基于64个可打印字符来表示二进制数据的表示方法。
用途：
1）所有的二进制文件，都可以因此转化为可打印的文本编码（都变成ASCII码可打印字符），使用文本软件进行编辑。
2）能够对文本进行简单的加密。

## Base64编码规则

例子：
Text content      M          a           n
ASCII             77         97          110
Bit pattern     01001101     01100001    01101111(每六个bite为一组，不够的在前面补充00)
Index            19       22          5        46
Base64-Encoded   T        W           F         u

## Base32和Base16

原理和Base64的原理一样，只是他们可以转换的字符越来越少了。