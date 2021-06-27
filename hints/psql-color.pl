#!/usr/bin/perl -w

$mode = 0;
while (<>) {
	if ($. == 1 && /^-\[/) {
		$mode = 1;
	}
	if ($. == 1 && /^\s+(Table "|List of )/) {
		$mode = 2;
	}
	if ($mode == 0) { # table
		if ($. == 1) {
			s/^(.*)$/\033[1;33m$1\033[0m/;
			s/ \| / \033[1;32m|\033[1;33m /g;
		} elsif ($. == 2) {
			s/^([-+]+)$/\033[1;32m$1\033[0m/;
		} else {
			s/ \| / \033[1;32m|\033[0m /g;
			s/(¤)/\033[1;31m$1\033[0m/g;
		}
	} elsif ($mode == 1) { # rows table
		s/^(-\[ RECORD \d+ \][-+]+)$/\033[1;32m$1\033[0m/;
		s/ \| / \033[1;32m|\033[0m /g;
		s/(¤)/\033[1;31m$1\033[0m/g;
	} else { # \d
		if ($. == 1) {
			s/^(.*)$/\033[1;37m$1\033[0m/;
		} elsif ($. == 2) {
			s/^(.*)$/\033[1;33m$1\033[0m/;
			s/ \| / \033[1;32m|\033[1;33m /g;
		} elsif ($. == 3) {
			s/^([-+]+)$/\033[1;32m$1\033[0m/;
		} else {
			s/ \| / \033[1;32m|\033[0m /g;
			s/^(Indexes:)$/\033[1;33m$1\033[0m/;
		}
	}
	print;
}
