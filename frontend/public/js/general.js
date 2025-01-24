const checkLogin = () => {
  if (localStorage.getItem('token')) {
    return true
  }
  return false
}
