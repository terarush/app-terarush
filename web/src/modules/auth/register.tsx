import { AuthLayout } from "./layouts/AuthLayout"
import { RegisterForm } from "./components/RegisterForm"

export default function Register() {
  return (
    <AuthLayout>
      <RegisterForm />
    </AuthLayout>
  )
}
