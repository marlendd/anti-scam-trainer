import { RegisterForm } from '@/features/auth'
import { useDocumentTitle } from '@/shared/lib/use-document-title'

export function RegisterPage() {
    useDocumentTitle('Регистрация')

    return <RegisterForm />
}
