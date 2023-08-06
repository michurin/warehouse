#!/usr/bin/perl -w

sub gx {
	local $_ = $_[0];
	my $m = $_[1] || '';
	if (s/\[ RECORD (\d+) \]/━━━━━ $1 ━━━━━/) {
		s/[-+]{10,}/'━' x 20/ge;
		s/[-+]/━/g;
	}
	s/\-/─/g;
	s/\+/┼/g;
	s/\|/│/g;
	return "$m\033[92m${_}\033[0m$m";
}

$mode = 0;
while (<>) {
	if ($. == 1 && /^-\[/) {
		$mode = 1;
	}
	if ($. == 1 && /^\s*(Table "|Index "|View "|Sequence "|List of )/) {
		$mode = 2;
	}
	if ($mode == 0) { # table
		if ($. == 1) {
			s/([\w\d]+)/\033[1;33m$1\033[0m/g;
			s/ (\|) /gx($1, ' ')/eg;
		} elsif ($. == 2) {
			s/^([-+]+)$/gx($1)/e;
		} else {
			s/ (\|) /gx($1, ' ')/ge;
			s/(¤)/\033[1;31m$1\033[0m/g;
		}
	} elsif ($mode == 1) { # rows table
		s/^(-\[ RECORD \d+ \][-+]+)$/gx($1)/e;
		s/ (\|) /gx($1, ' ')/ge;
		s/(¤)/\033[1;31m$1\033[0m/g;
	} elsif ($mode == 2) { # \d
		if ($. == 1) {
			s/^(.*)$/\033[1;37m$1\033[0m/;
		} elsif ($. == 2) {
			s/^(.*)$/\033[1;33m$1\033[0m/;
			s/ (\|) /gx($1, ' ')/ge;
		} elsif ($. == 3) {
			s/^([-+]+)$/gx($1)/e;
		} else {
			if (!s/ (\|) /gx($1, ' ')/ge) {
				if (s/^((Indexes|Triggers):)$/\033[1;33m$1\033[0m/) {
					$mode = 3;
				}
			}
		}
	} else { # \dS trailing list
		s/\(([^)]+)\)/(\033[1;33m$1\033[0m)/;
		s/"([^"]+)"/"\033[92m$1\033[0m"/;
		s/^((Indexes|Triggers):)$/\033[1;33m$1\033[0m/;
	}
	print;
}
