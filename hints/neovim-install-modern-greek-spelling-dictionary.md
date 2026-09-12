# Το πρόβλημα

-σ -> -ς

# Solution

```
git clone https://github.com/bataak/dictionaries-hunspell.git
cd dictionaries-hunspell
cp dictionaries/el/index.aff el.aff
cp dictionaries/el/index.dic el.dic
nvim -u NONE +'mkspell! ~/.local/share/nvim-12/site/spell/el ./el.dic' +q
```

(check path to your dictionaries)
