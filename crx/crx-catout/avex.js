// avito link in ms exchange
(() => {
  const F = () => {
    Array.from(document.querySelectorAll('*'))
      .map(x => x.childNodes.length === 1 ? (x.childNodes[0].nodeType === 3 ? (x.innerText.match('^http') ? x : false) : false) : false)
      .filter(x => x)
      .forEach(x => {
      c=x.childNodes[0];
      t=x.innerText;
      x.removeChild(c);
      a=document.createElement('a');
      a.href=t;
      a.innerText=t.replace(/[^/]*:\/\//, '');
      a.target='_blank';
      a.style.color='#080';
      a.style.fontWeight='bold';
      x.append(a);
    });
  };
  setTimeout(F, 500);
  setTimeout(F, 1000);
  setInterval(F, 2000);
})();
