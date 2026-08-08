import { useNavigate } from 'react-router-dom'

import { WelcomeSlider } from '@/widgets/welcome-slider'

export function WelcomePage() {
    const navigate = useNavigate()

    function handleComplete() {
        navigate('/training/role-selection')
    }

    return (
        <WelcomeSlider
            onComplete={handleComplete}
            onSkip={handleComplete}
        />
    )
}