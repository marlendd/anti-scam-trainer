import { LoginForm } from '@/features/auth'
import { useDocumentTitle } from '@/shared/lib/use-document-title'

export function LoginPage() {
  useDocumentTitle('Вход')

  return <LoginForm />
}