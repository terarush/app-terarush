import { AuthLayout } from "./layouts/AuthLayout"
import { LoginForm } from "./components/LoginForm"
import { authContent } from "./content/auth"

export default function Login() {
  return (
    <AuthLayout title={authContent.login.title} subtitle={authContent.login.subtitle}>
      <LoginForm />
    </AuthLayout>
  )
}
