function before() {
    return true
}

// #region empty
/* begin */
/* end */
// #endregion

// #region nonempty
/* begin */
function nonempty() {
    return "Hello"
}
/* end */
// #endregion

/* #region block */
/* begin */
function block() {
    return true
}
/* end */
/* #endregion */

function after() {
    return false
}
