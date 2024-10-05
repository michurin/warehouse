# Linux on MacOS, Arch Linux, Debian Linux, sharing folders

## Arch Linux

- qemu
- host: MacOS
- guest: Arch Linux

```sh
create -f qcow2 image.img 100G

curl http://ftp.agdsn.de/pub/mirrors/archlinux/iso/2023.06.01/archlinux-2023.06.01-x86_64.iso -o boot.iso
curl http://ftp.agdsn.de/pub/mirrors/archlinux/iso/2023.06.01/arch/boot/x86_64/initramfs-linux.img -o initramfs-linux.img
curl http://ftp.agdsn.de/pub/mirrors/archlinux/iso/2023.06.01/arch/boot/x86_64/vmlinuz-linux -o vmlinuz-linux

qemu-system-x86_64 -m 2048 -vga virtio -show-cursor -usb -device usb-tablet -enable-kvm -cdrom boot.iso -drive file=image.img,if=virtio -accel hvf -cpu host -boot d

fdisk -l
fdist /dev/vda
mkfs.ext4 /dev/vda1
mount /dev/vda1 /mnt
pacstrap -K /mnt base linux linux-firmware mc tmux
genfstab -U /mnt >> /mnt/etc/fstab
arch-chroot /mnt
mkinitcpio -P
passwd
pacman -Suy grub
grub-install --target=i386-pc /dev/vda
grub-mkconfig -o /boot/grub/grub.cfg
pacman -Suy netctl dialog dhclient dhcpcd
pacman -Suy openssh # PermitRootLogin yes
pacman -Suy neovim
pacman -Suy sudo
pacman -Suy inetutils # <-telnet
exit
umount /mnt
halt -p

qemu-system-x86_64 -m 2048 -vga virtio -hda image.img -accel hvf -cpu host -boot c -net user,hostfwd=tcp::10022-:22 -net nic

dhclient ens3
vi /etc/netctl/ens3
netclt enable ens3

ssh-keygen -A
systemctl enable sshd

useradd -m a
passwd a
su a
```

## sshFS

```sh
sshfs macuser@10.0.2.2:~ ~/tmp
```

## Debian (with shared FS)

```sh
curl https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.7.0-amd64-netinst.iso -o debian-12.7.0-amd64-netinst.iso -L
qemu-img create -f qcow2 image.qcow2 10G
qemu-system-x86_64 -drive file=image.qcow2,format=qcow2 -m 1G -accel kvm -cdrom debian-12.7.0-amd64-netinst.iso -boot order=d

qemu-system-x86_64 -drive file=image.qcow2,format=qcow2 -m 1G -accel kvm -nic user,hostfwd=tcp::2222-:22 -virtfs local,path=$HOME/shared,mount_tag=shared,security_model=mapped-xattr
```

```sh
mount -t 9p shared /mnt
mount -t 9p -o trans=virtio,version=9p2000.L shared /mnt/shared
```

more details `https://wiki.qemu.org/Documentation/9psetup`
