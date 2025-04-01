## notes:

- since wpm is the standard we evaluate word by word!
- logical worde devider is `<space>`
- overtype (ot) is when type letters > current word letter
- undertype is when word is not overtyped (assuing player is dont typing current word)
- space devider edgecase : undertyped and space is pressed -> the standard is just jump to the next word and mark current one as incorrect (I personally don't like it space should count as a letter and it's easier to implement but we go for the standard!)
- the placeholder should shift forward in overtype and if we assign `<space>` as the begining of the next word while countong it as a logical devider we don't have to deal with problems caused by overtype shift in validation! another solution is to refactor validation by having and expected next word conect, and since we know where spaces are, we shift if we get letter in space place and we jump to next word if we get space in letter's place! the rest is even easier to validate!

## dev todos:

- [ ] find a way live render the the paragraph as placeholder and impliment the eval function in it!
- [ ] add a test for the new `--no-stdin` flag
- [ ] probaly a good time to start testing for different shells and terminal emulators!
- [ ] turns out it's really hard to deal with the placeholder shift in the lip textarea! best case i was getting ANSI chars in the nstdout! solution: rawdog a textarea with hanging placeholder, for now i just changed the gameplay by shifting typed words up and resetting the textarea! (can be useful for simon says gamemode)

## gameplay todos:

- [ ] check out typing styles! now everyone uses vim mode, eg: cltr + backspace should delete the whole word, (and we should support visual mode, diw)
- [ ] lowercase , no digit as default gamemode ?

## server todos

- [ ] we can use the vercel runtime for free hosting if the server doesnt get complicated
- [ ] sources for text : wikipedia, ...
