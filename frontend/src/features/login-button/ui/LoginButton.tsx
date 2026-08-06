import {Button} from "@/shared/ui/button";
import {useNavigate} from "react-router-dom";
export function LoginButton() {

    const navigate = useNavigate();

    const handleClick = () => {
        navigate('/login');
    }

    return (
        <Button size='small' variant='secondary' onClick={handleClick}>
            Войти
        </Button>
    )
}