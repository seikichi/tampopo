#!/bin/bash -e

DIR=$(dirname $0)/LOVE_HINA_01

mkdir -p $DIR
wget -P $DIR -O LH01_hq.pdf "http://dl.j-comi.jp/download/book/101/hq/LH01_hq.pdf"
gs -q -dNOPAUSE -dBATCH -sDEVICE=pdfwrite -sOutputFile=${DIR}/LH01_hq_decrypted.pdf -c .setpdfwrite -f $DIR/LH01_hq.pdf
# pdfimages -j ${DIR}/LH01_hq_decrypted.pdf ${DIR}/page
