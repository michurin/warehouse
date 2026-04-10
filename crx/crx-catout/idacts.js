// github: highlight users
(() => {
  const users = {
    'sektor-miras': { color: '#0f0' },
    'dauletsailauov': { color: '#f00' },
    'alexeymichurin': { color: '#fff' },
  }
  Object.values(users).forEach((u) => {
    u.style = 'border: 1px dotted ' + u.color + '; border-radius: .3em; margin: -1px; color: ' + u.color + ' !important'
  })
  const hl = () => {
    if (!location.pathname.endsWith('/actions')) {
      console.log('skip page')
      return
    }
    document.querySelectorAll('[data-hovercard-type="user"]').forEach((e) => {
      const user = e.innerText
      if (users[user]) {
        e.style = users[user].style
      }
    })
  };
  hl()
  setTimeout(hl, 200)
  setTimeout(hl, 500)
  setInterval(hl, 1000)
})();
