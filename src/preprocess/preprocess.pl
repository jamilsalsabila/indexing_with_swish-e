#!/usr/bin/perl -w
use strict;

# preprocess(inputFile, outputFile)
sub preprocess {
    # open file
    open(INPUT, "<:encoding(UTF-8)", $_[0]) or die $!;
    open(OUTPUT, ">:encoding(UTF-8)", $_[1]) or die $!;

    # read file line by line: search for lines containing "content: askaks"
    my @lines = <INPUT>;
    my $i = 0;
    my $string = "";
    # print @lines;
    while ($i < scalar @lines) {
        if ($lines[$i] =~ m/^content:/i) {
            $lines[$i] =~ s/^content:\s*//i;
            last;
        }
        $i++;
    }

    while ($i < scalar @lines) {
        $string .= $lines[$i];
        $i++
    }

    $string =~ s/[\n\r]/ /gi;
    $string =~ s/-/ /gi;
    $string =~ s/[“,:]//gi;
    $string =~ s/([\w\W])\.(\s?)/$1$2/gi;
    $string =~ s/[^A-Za-z0-9\s)(\/.]//gi;
    $string =~ s/\s+/ /gi;
    $string =~ s/^ //gi;
    $string = lc($string);

    print OUTPUT $string;

    # close file 
    close INPUT;
    close OUTPUT;
}

my $outputPath = "ALL/";
my @inputPath = ("antara/entertainment/", "antara/oto/", "antara/tekno/", "antara/sport/", "kompas/entertainment/", "kompas/oto/", "kompas/tekno/", "kompas/sport/", "tempo/entertainment/", "tempo/oto/", "tempo/tekno/", "tempo/sport/");
#

for (my $p = 0; $p < scalar @inputPath; $p++) {
    for (my $i = 0; $i < 10; $i++) {
        opendir DIR, "$inputPath[$p]$i/" or die $!;
        my @folder = split("/", "$inputPath[$p]");
        while (my $file = readdir DIR) {
            if ($file eq "." or $file eq "..") {
                next;
            }
            print "Processing: $inputPath[$p]$i/$file\n";
            preprocess("$inputPath[$p]$i/$file", "$outputPath$folder[0]/$file-$folder[1].txt"); 
        } 
        close DIR;  
    }
}

