(() => {
  console.log('HERE!');
  const f = () => {
    x = document.getElementsByTagName('div')
    for (i = 0; i < x.length; i++) {
      if (x[i].innerText === 'Не сейчас') {
        x[i].click() // close annoying popup "share your geo position" (ru)
        console.log('HERE: Click!');
      }
    }
  };
  for (i = 500; i < 10000; i += 500) {
    setTimeout(f, i);
  }
})();
