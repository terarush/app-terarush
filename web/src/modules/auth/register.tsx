import { AuthLayout } from "./layouts/AuthLayout"
import { RegisterForm } from "./components/RegisterForm"
import { authContent } from "./content/auth"

export default function Register() {
  return (
    <AuthLayout title={authContent.register.title} subtitle={authContent.register.subtitle}>
      <RegisterForm />
    </AuthLayout>
  )
}
