// github: highlight users
(() => {
  const users = [
    { nik: 'sektor-miras', color: '#0f0' },
    { nik: 'dauletsailauov', color: '#f00' },
    { nik: 'alexeymichurin', color: '#fff' },
  ]
  users.forEach((u) => {
    /*
    u.style = 'text-shadow: ' +
      '1px 1px 1px ' + u.color + ', ' +
      '1px 0px 1px ' + u.color + ', ' +
      '1px -1px 1px ' + u.color + ', ' +
      '0px 1px 1px ' + u.color + ', ' +
      '0px 0px 1px ' + u.color + ', ' +
      '0px -1px 1px ' + u.color + ', ' +
      '-1px 1px 1px ' + u.color + ', ' +
      '-1px 0px 1px ' + u.color + ', ' +
      '-1px -1px 1px ' + u.color + '; ' +
      'border: 1px dotted ' + u.color + '; margin: -1px; ' +
      'color: black !important'
    */
    u.style = 'border: 1px dotted ' + u.color + '; border-radius: .3em; margin: -1px; color: ' + u.color + ' !important';
  })
  const hl = () => {
    if (!location.pathname.endsWith('/pulls')) {
      console.log('skip page')
      return
    }
    users.forEach((u) => {
      document.querySelectorAll('a[data-hovercard-url="/users/' + u.nik + '/hovercard"]').forEach((e) => {
        e.style = u.style
      })
    })
    document.querySelectorAll('a[aria-label="Review required before merging"]').forEach((e) => {
      e.style = 'border-width: 0px; border-radius: .3em; background-color: #800; color: #000 !important'
    })
  }
  hl()
  setTimeout(hl, 200)
  setTimeout(hl, 500)
  setTimeout(hl, 1000)
  setTimeout(hl, 2000)
  setInterval(hl, 5000)
})();
