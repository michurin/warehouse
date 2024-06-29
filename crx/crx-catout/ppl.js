(() => {
  const work = async () => {
    console.log('PPL RUN...')

    const p = document.querySelectorAll('h2')
    console.log(p.length)
    if (p.length !== 1) {
      return
    }

    const l = location.pathname.split('/')
    if (l.length <= 0) {
      return
    }
    const id = l[l.length - 1]

    const x = await fetch('https://people-api.avito.ru/service-proxy/ap-profile-composition/web/v1/profile?user_id=' + id, {
      credentials: 'include',
      headers: { 'Accept': 'application/json' }
    })
    const y = await x.json()
    const name = y.middleName

    const e = document.createElement('span')
    e.innerText = ' [' + name + ']'
    p[0].append(e)
  }

  if (document.readyState == 'loading') {
    // still loading, wait for the event
    document.addEventListener('DOMContentLoaded', work);
  } else {
    // DOM is ready!
    work();
  }
})();
