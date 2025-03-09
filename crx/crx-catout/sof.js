(() => {
  console.log('HERE!');
  const f = () => {
    const e = document.getElementById('onetrust-accept-btn-handler');
    if (e) {
      e.click();
    }
  };
  for (i = 500; i < 10000; i += 500) {
    setTimeout(f, i);
  }
})();
