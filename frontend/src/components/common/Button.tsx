import { ButtonHTMLAttributes } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary';
}

function Button({ variant = 'primary', className = '', ...props }: ButtonProps) {
  const baseStyles = 'px-6 py-4 rounded-2xl font-medium transition-all duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-transparent';
  const variants = {
    primary: 'text-white shadow-lg shadow-violet-500/25 hover:shadow-violet-500/40',
    secondary: 'text-gray-600 hover:text-gray-800 shadow-sm',
  };

  return (
    <button
      className={`${baseStyles} ${variants[variant]} ${className}`}
      {...props}
    />
  );
}

export default Button; 