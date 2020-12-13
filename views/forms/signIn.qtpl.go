// Code generated by qtc from "signIn.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// All the text outside function templates is treated as comments,
// i.e. it is just ignored by quicktemplate compiler (`qtc`). It is for humans.
//
// .

//line views/forms/signIn.qtpl:5
package forms

//line views/forms/signIn.qtpl:5
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/forms/signIn.qtpl:5
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/forms/signIn.qtpl:5
func StreamSignInForm(qw422016 *qt422016.Writer) {
//line views/forms/signIn.qtpl:5
	qw422016.N().S(`
<form target="content" action="/api/v1/user/signin/" class="form-signin" method="post"  enctype="multipart/form-data"
 onsubmit="return saveForm(this, afterLogin);" novalidate>
        <h2 class="form-signin-heading">Авторизация</h2>
        <input type="hidden" name="url" value=""/>
        <input type="email" name="email" class="input-block-level" placeholder="Email, указанный при регистрации" value="zero@null.com">
        <input type="password" name="key" class="input-block-level" placeholder="Введите пароль, полученный по почте">
        <label class="checkbox">
         <input type="checkbox" name="remember" value="remember-me"> Запомнить меня в системе
        </label>
       <button class="main-btn" type="submit">Войти</button>
     <output></output>
     <progress value='0' max='100' hidden > </progress>
</form>
`)
//line views/forms/signIn.qtpl:19
}

//line views/forms/signIn.qtpl:19
func WriteSignInForm(qq422016 qtio422016.Writer) {
//line views/forms/signIn.qtpl:19
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/forms/signIn.qtpl:19
	StreamSignInForm(qw422016)
//line views/forms/signIn.qtpl:19
	qt422016.ReleaseWriter(qw422016)
//line views/forms/signIn.qtpl:19
}

//line views/forms/signIn.qtpl:19
func SignInForm() string {
//line views/forms/signIn.qtpl:19
	qb422016 := qt422016.AcquireByteBuffer()
//line views/forms/signIn.qtpl:19
	WriteSignInForm(qb422016)
//line views/forms/signIn.qtpl:19
	qs422016 := string(qb422016.B)
//line views/forms/signIn.qtpl:19
	qt422016.ReleaseByteBuffer(qb422016)
//line views/forms/signIn.qtpl:19
	return qs422016
//line views/forms/signIn.qtpl:19
}
