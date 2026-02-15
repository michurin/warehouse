// github: highlight users
(() => {
  const users = {
    'sektor-miras': { color: '#0f0' },
    'dauletsailauov': { color: '#f00' },
    'alexeymichurin': { color: '#fff' },
  }
  Object.values(users).forEach((u) => {
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
    document.querySelectorAll('.js-issue-row').forEach((sec) => {
      var type = 'none'
      var user = ''
      var userElement
      var labelElement
      var titleElement
      sec.querySelectorAll('.Link--muted').forEach((e) => {
        var text = e.innerText
        if (text == 'Draft') {
          type = 'draft'
          labelElement = e
          titleElement = sec.querySelector('.markdown-title')
        }
        if (text == 'Approved') {
          type = 'approved'
          labelElement = e
        }
        if (e.dataset.hovercardType == 'user') {
          user = text
          userElement = e
        }
      })
      if (type == 'draft') {
        labelElement.style = 'border-width: 0px; border-radius: .3em; background-color: #444; color: #000 !important'
        sec.style = 'background: repeating-linear-gradient(-40deg, rgba(0,0,0,0) 0px, rgba(0,0,0,0) 3px, rgba(255,255,255,.1) 5px, rgba(0,0,0,0) 7px, rgba(0,0,0,0) 10px)'
        //titleElement.style = 'color: #000 !important; text-shadow: 0px 0px 8px #ffffff;'
        titleElement.style = 'color: #333 !important; -webkit-text-stroke: 1px #777 !important;'
        if (user != 'alexeymichurin') {
          return
        }
      }
      if (type == 'approved') {
        labelElement.style = 'border-width: 0px; border-radius: .3em; background-color: #080; color: #000 !important'
      }
      if (users[user]) {
        userElement.style = users[user].style
        if (type == 'approved') {
          sec.style = 'background: repeating-linear-gradient(-40deg, rgba(0,0,0,0) 0px, rgba(0,0,0,0) 3px, rgba(0,255,0,.3) 5px, rgba(0,0,0,0) 7px, rgba(0,0,0,0) 10px)'
        }
      }
    })
  }
  hl()
  setTimeout(hl, 200)
  setTimeout(hl, 500)
  setInterval(hl, 1000)
})();
