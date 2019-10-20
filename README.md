# Unique files

This is a litle script to find duplicate files in various folders.

This is usefull when importing many files from many sources, and we don't know which ones are identical.

## Usage

```bash
  go run main.go --path path/to/folder/1 --path path/to/folder/2
```

## Example use case:

After an event, many people share their pictures. But then you get duplicates because A shared his files with you, so did B, but B had already received A's photos when he shared, so now A's photos are duplicates
