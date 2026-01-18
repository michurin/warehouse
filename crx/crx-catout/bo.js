// boundary
(() => {
  console.log("bo: Start")
  setInterval(() => {
    console.log("bo: Check")
    if (window.location.pathname == '/authentication-complete') {
      window.close()
    }
  }, 1000)
})()
