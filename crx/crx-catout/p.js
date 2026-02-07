// /pogoda
(() => {
  setInterval(() => {
    if (!location.pathname.includes('/pogoda/')) {
      console.log('skip page')
      return
    }
    document.querySelectorAll('section').forEach(e => {
      if (e.className.startsWith('AppWidgetNowcast_container_logo_')) { e.remove() }
      if (e.className.startsWith('Money_ecomFooter_')) { e.remove() }
    })
    document.querySelectorAll('div').forEach(e => {
      if (e.className.startsWith('MainPage_topBlockWithMoney_')) { e.remove() }
      // if (e.className.startsWith('AppPromoInner_container_')) { e.parentNode.remove() } // too aggressive?
    })
    document.querySelectorAll('li').forEach(e => {
      if (e.className.startsWith('AppForecastMoney_wrap_')) { e.remove() }
    })
    document.querySelectorAll('aside').forEach(e => {
      if (e.className.startsWith('AppLayoutTypeMain_right_')) { e.remove() }
    })
  }, 2000)
})()
