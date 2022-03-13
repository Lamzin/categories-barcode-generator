To generate new barcodes do next:
1. Download Google spreadsheet with categories as CVS file.
2. Copy CSV file content to `data/category_list.csv`
3. Run `docker run -v "$(pwd)"/out:/out --rm -it $(docker build -q .)`
4. Find generated PDF files in `out/` folder.
