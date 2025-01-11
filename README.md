# Kindle HTTP Server

This is a simple http server that can be used to remotely control a kindle (most
importantly -- to turn pages).

Tested on Kindle Paperwhite 3rd gen.

# Kindle preparation

1. Jailbreak the kindle
2. Install MR Package Installer and KUAL
3. Set up networking and ssh access

# Building and installing

Build, install, and run server on the kindle using the following command:

```sh
make KINDLE_SSH_ADDRESS=root@kindle.local all install run-background
```

This will build the binary, copy it to `/tmp` and start it in the background. I
think it wouldn't be too hard to set it up to run at startup using an
`/etc/upstart` script (see the resources). But I didn't want to mess with the
boot sequence and potentially break it. I just start the server manually after a
restart, which doesn't happen often.

# Usage

```sh
$ curl http://kindle.local/lipc-get-prop/com.lab126.amazonRegistrationService/GetDeviceName
"Some Kindle"
$ curl http://kindle.local/battery-level
58
$ curl http://kindle.local/swipe-left
"OK"
$ curl http://kindle.local/swipe-right
"OK"
$ curl http://kindle.local/toggle-power-button
```

# Resources

- [How to Jailbreak the Kindle Scribe Using the LanguageBreak Exploit - Noah Nash](https://www.noahnash.net/blog/jailbreak-kindle-scribe/)
- [LanguageBreak - Jailbreak for any kindle running FW 5.16.2.1.1 or LOWER](https://github.com/notmarek/LanguageBreak)
- [SSHing to My Kindle - yaleman.org](https://yaleman.org/post/2018/2018-08-06-sshing-to-my-kindle/)
- [USBnetwork SSH(terminal) on Windows 10:How I did it - MobileRead Forum](https://www.mobileread.com/forums/showthread.php?t=340208)
- [Best way to exec script on boot - MobileRead Forum](https://www.mobileread.com/forums/showthread.php?t=221019)
- [Kindle hacking: a deeper dive into the internals - SixFoisNeuf](https://www.sixfoisneuf.fr/posts/kindle-hacking-deeper-dive-internals/)
- [Kindle Apps & Services - Kindle Modding Wiki](https://kindlemodding.org/kindle-apps-and-services/)
- [Kindle Hacks Information - MobileRead Wiki](https://wiki.mobileread.com/wiki/Kindle_Hacks_Information)
- [Kindle Touch Hacking - MobileRead Wiki](https://wiki.mobileread.com/wiki/Kindle_Touch_Hacking)
