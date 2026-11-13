package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSplitStringIntoChunks_Short(t *testing.T) {
	if got := splitStringIntoChunks(""); len(got) != 0 {
		t.Fatalf("expected no chunks, got %d", len(got))
	}
	got := splitStringIntoChunks("Hello world.")
	if len(got) != 1 || got[0] != "Hello world." {
		t.Fatalf("unexpected chunks: %q", got)
	}
}

func TestSplitStringIntoChunks_LongWordWithoutSeparator(t *testing.T) {
	input := strings.Repeat("ż", 5000)
	chunks := splitStringIntoChunks(input)
	for i, c := range chunks {
		if len(c) > 4000 {
			t.Errorf("chunk %d too long: %d", i, len(c))
		}
		if !utf8.ValidString(c) {
			t.Errorf("chunk %d is not valid UTF-8", i)
		}
	}
	if strings.Join(chunks, "") != input {
		t.Error("joined chunks differ from input")
	}
}

func TestSplitStringIntoChunks_LongText(t *testing.T) {
	input := `Date: Thu, 12 Nov 2026 06:46:17 +0000
Subject: [reboot required] unattended-upgrades result for XXXXXXX: SUCCESS
From: root
To: user@example.com
Auto-Submitted: auto-generated
MIME-Version: 1.0
Content-Type: text/plain; charset="utf-8"
Content-Transfer-Encoding: quoted-printable

Unattended upgrade result: All upgrades installed

Warning: A reboot is required to complete this upgrade, or a previous one.

Packages that were upgraded:
 apt apt-transport-https apt-utils bind9-dnsutils bind9-host bind9-libs
 ca-certificates distro-info distro-info-data fluent-bit ghostscript
 git git-man gitlab-runner gitlab-runner-helper-images
 gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-x
 libapt-pkg6.0 libbrotli1 libglib2.0-0 libglib2.0-bin libglib2.0-data
 libgs9 libgs9-common libgstreamer-plugins-base1.0-0
 libgstreamer-plugins-good1.0-0 liblzma5 libmspack0 libmysqlclient21
 libnss3 libpcre2-8-0 libperl5.34 libsoup2.4-1 libsoup2.4-common
 libxtst6 linux-headers-generic linux-headers-virtual
 linux-image-virtual linux-virtual nano perl perl-base
 perl-modules-5.34 python3-distro-info python3-jwt python3-pyasn1
 python3-twisted snapd tmux ubuntu-advantage-tools ubuntu-pro-client
 ubuntu-pro-client-l10n xfsprogs xz-utils

Packages that were auto-removed:
 linux-headers-5.15.0-194 linux-headers-5.15.0-194-generic
 linux-modules-5.15.0-194-generic linux-image-5.15.0-194-generic

Package installation log:
Log started: 2026-11-12  06:36:30
Preparing to unpack .../libnss3_2%3a3.98-0ubuntu0.22.04.5_amd64.deb ...
Unpacking libnss3:amd64 (2:3.98-0ubuntu0.22.04.5) over (2:3.98-0ubuntu0.22.=
04.4) ...
Setting up libnss3:amd64 (2:3.98-0ubuntu0.22.04.5) ...
Processing triggers for libc-bin (2.35-0ubuntu3.15) ...

Running kernel seems to be up-to-date.

No services need to be restarted.

No containers need to be restarted.

No user sessions are running outdated binaries.

No VM guests are running outdated hypervisor (qemu) binaries on this host.
Log ended: 2026-11-12  06:36:35

Log started: 2026-11-12  06:36:36
Selecting previously unselected package linux-headers-5.15.0-203.
Preparing to unpack .../0-linux-headers-5.15.0-203_5.15.0-203.213_all.deb .=
..
Unpacking linux-headers-5.15.0-203 (5.15.0-203.213) ...
Selecting previously unselected package linux-headers-5.15.0-203-generic.
Preparing to unpack .../1-linux-headers-5.15.0-203-generic_5.15.0-203.213_a=
md64.deb ...
Unpacking linux-headers-5.15.0-203-generic (5.15.0-203.213) ...
Selecting previously unselected package linux-modules-5.15.0-203-generic.
Preparing to unpack .../2-linux-modules-5.15.0-203-generic_5.15.0-203.213_a=
md64.deb ...
Unpacking linux-modules-5.15.0-203-generic (5.15.0-203.213) ...
Selecting previously unselected package linux-image-5.15.0-203-generic.
Preparing to unpack .../3-linux-image-5.15.0-203-generic_5.15.0-203.213_amd=
64.deb ...
Unpacking linux-image-5.15.0-203-generic (5.15.0-203.213) ...
Preparing to unpack .../4-linux-virtual_5.15.0.203.175_amd64.deb ...
Unpacking linux-virtual (5.15.0.203.175) over (5.15.0.198.171) ...
Preparing to unpack .../5-linux-image-virtual_5.15.0.203.175_amd64.deb ...
Unpacking linux-image-virtual (5.15.0.203.175) over (5.15.0.198.171) ...
Preparing to unpack .../6-linux-headers-virtual_5.15.0.203.175_amd64.deb ...
Unpacking linux-headers-virtual (5.15.0.203.175) over (5.15.0.198.171) ...
Preparing to unpack .../7-linux-headers-generic_5.15.0.203.175_amd64.deb ...
Unpacking linux-headers-generic (5.15.0.203.175) over (5.15.0.198.171) ...
Setting up linux-headers-5.15.0-203 (5.15.0-203.213) ...
Setting up linux-headers-5.15.0-203-generic (5.15.0-203.213) ...
Setting up linux-headers-generic (5.15.0.203.175) ...
Setting up linux-headers-virtual (5.15.0.203.175) ...
Setting up linux-image-5.15.0-203-generic (5.15.0-203.213) ...
I: /boot/vmlinuz.old is now a symlink to vmlinuz-5.15.0-198-generic
I: /boot/initrd.img.old is now a symlink to initrd.img-5.15.0-198-generic
I: /boot/vmlinuz is now a symlink to vmlinuz-5.15.0-203-generic
I: /boot/initrd.img is now a symlink to initrd.img-5.15.0-203-generic
Setting up linux-image-virtual (5.15.0.203.175) ...
Setting up linux-modules-5.15.0-203-generic (5.15.0-203.213) ...
Setting up linux-virtual (5.15.0.203.175) ...
Processing triggers for linux-image-5.15.0-203-generic (5.15.0-203.213) ...
/etc/kernel/postinst.d/initramfs-tools:
update-initramfs: Generating /boot/initrd.img-5.15.0-203-generic
/etc/kernel/postinst.d/zz-update-grub:
Sourcing file /etc/default/grub
Sourcing file /etc/default/grub.d/50-cloudimg-settings.cfg
Sourcing file /etc/default/grub.d/init-select.cfg
Generating grub configuration file ...
Found linux image: /boot/vmlinuz-5.15.0-203-generic
Found initrd image: /boot/initrd.img-5.15.0-203-generic
Found linux image: /boot/vmlinuz-5.15.0-198-generic
Found initrd image: /boot/initrd.img-5.15.0-198-generic
Found linux image: /boot/vmlinuz-5.15.0-194-generic
Found initrd image: /boot/initrd.img-5.15.0-194-generic
Warning: os-prober will not be executed to detect other bootable partitions.
Systems on them will not be added to the GRUB boot configuration.
Check GRUB_DISABLE_OS_PROBER documentation entry.
done

Pending kernel upgrade!

Running kernel version:
5.15.0-198-generic

Diagnostics:
The currently running kernel version is not the expected kernel version 5=
.15.0-203-generic.

Restarting the system to load the new kernel will not be handled automatica=
lly, so you should consider rebooting. [Return]

No services need to be restarted.

No containers need to be restarted.

No user sessions are running outdated binaries.

No VM guests are running outdated hypervisor (qemu) binaries on this host.
Log ended: 2026-11-12  06:36:54

Log started: 2026-11-12  06:36:54
Preparing to unpack .../ubuntu-pro-client-l10n_38ubuntu0~22.04_amd64.deb ...
Unpacking ubuntu-pro-client-l10n (38ubuntu0~22.04) over (37.2ubuntu~22.04.1=
) ...
Preparing to unpack .../ubuntu-pro-client_38ubuntu0~22.04_amd64.deb ...
Unpacking ubuntu-pro-client (38ubuntu0~22.04) over (37.2ubuntu~22.04.1) ...
Setting up ubuntu-pro-client (38ubuntu0~22.04) ...
Installing new version of config file /etc/apparmor.d/ubuntu_pro_esm_cache =
...
Setting up ubuntu-pro-client-l10n (38ubuntu0~22.04) ...
Processing triggers for man-db (2.10.2-1) ...

Pending kernel upgrade!

Running kernel version:
5.15.0-198-generic

Diagnostics:
The currently running kernel version is not the expected kernel version 5=
.15.0-203-generic.

Restarting the system to load the new kernel will not be handled automatica=
lly, so you should consider rebooting. [Return]

No services need to be restarted.

No containers need to be restarted.

No user sessions are running outdated binaries.

No VM guests are running outdated hypervisor (qemu) binaries on this host.
Log ended: 2026-11-12  06:37:00

Log started: 2026-11-12  06:37:01
Preparing to unpack .../libglib2.0-bin_2.72.4-0ubuntu2.11_amd64.deb ...
Unpacking libglib2.0-bin (2.72.4-0ubuntu2.11) over (2.72.4-0ubuntu2.10) ...
Preparing to unpack .../libglib2.0-0_2.72.4-0ubuntu2.11_amd64.deb ...
Unpacking libglib2.0-0:amd64 (2.72.4-0ubuntu2.11) over (2.72.4-0ubuntu2.10)=
...
Setting up libglib2.0-0:amd64 (2.72.4-0ubuntu2.11) ...
Setting up libglib2.0-bin (2.72.4-0ubuntu2.11) ...
Processing triggers for libc-bin (2.35-0ubuntu3.15) ...
Processing triggers for man-db (2.10.2-1) ...

Pending kernel upgrade!

Running kernel version:
5.15.0-198-generic

Diagnostics:
The currently running kernel version is not the expected kernel version 5=
.15.0-203-generic.

Restarting the system to load the new kernel will not be handled automatica=
lly, so you should consider rebooting. [Return]

Services to be restarted:
systemctl restart irqbalance.service
systemctl restart packagekit.service
systemctl restart polkit.service`

	chunks := splitStringIntoChunks(input)
	if len(chunks) != 3 {
		t.Fatalf("expected two chunks, got %d", len(chunks))
	}

	expectedEndOfChunk1 := `Setting up linux-headers-generic (5.15.0.203.175) ...
Setting up linux-headers-virtual (5.15.0.203.175) ...
Setting up linux-image-5.15.0-203-generic (5.15.0-203.213) ...`
	if !strings.HasSuffix(chunks[0], expectedEndOfChunk1) {
		t.Fatalf("expected chunk 1 to end with %q, got %q", expectedEndOfChunk1, chunks[0][len(chunks[0])-200:])
	}

	expectedStartOfChunk2 := `Subject: [reboot required] unattended-upgrades result for XXXXXXX: SUCCESS

I: /boot/vmlinuz.old is now a symlink to vmlinuz-5.15.0-198-generic
I: /boot/initrd.img.old is now a symlink to initrd.img-5.15.0-198-generic
I: /boot/vmlinuz is now a symlink to vmlinuz-5.15.0-203-generic`

	if !strings.HasPrefix(chunks[1], expectedStartOfChunk2) {
		t.Fatalf("expected chunk 2 to start with %q, got %q", expectedStartOfChunk2, chunks[1][0:200])
	}
}

func TestSplitMessageIntoChunks_RepeatsSubject(t *testing.T) {
	subject := "[reboot required] unattended-upgrades result for host: SUCCESS"
	body := strings.Repeat("alpha beta gamma delta.\n", 600)
	raw := "From: a@b.c\r\nSubject: " + subject + "\r\n\r\n" + body

	chunks := splitStringIntoChunks(raw)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for i, c := range chunks {
		if len(c) > maxChunkSize {
			t.Errorf("chunk %d too long: %d", i, len(c))
		}
		if i > 0 && !strings.HasPrefix(c, "Subject: "+subject+"\n") {
			t.Errorf("chunk %d does not start with subject: %q", i, c[:80])
		}
	}
}

func TestSplitMessageIntoChunks_EmptySubject(t *testing.T) {
	raw := "From: a@b.c\r\n\r\n" + strings.Repeat("alpha beta gamma delta.\n", 600)
	chunks := splitStringIntoChunks(raw)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for i, c := range chunks[1:] {
		if strings.HasPrefix(c, "Subject: \n") {
			t.Errorf("chunk %d does not start with empty subject: %q", i+1, c[:20])
		}
	}
}
