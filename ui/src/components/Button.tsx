import { ReactNode } from 'react';
import { classNames } from '.';
import React from 'react';

type HeroIcon = React.ForwardRefExoticComponent<
  React.PropsWithoutRef<React.SVGProps<SVGSVGElement>> & { title?: string; titleId?: string } & React.RefAttributes<SVGSVGElement>
>;

type ButtonVariants = 'primary' | 'secondary' | 'danger' | 'success' | 'warning';

interface Props extends React.DetailedHTMLProps<React.ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement> {
  children?: ReactNode;
  variant?: ButtonVariants;
  HeroIcon?: HeroIcon;
  loading?: boolean;
  type?: 'button' | 'submit' | 'reset';
}

const variants: { [key in ButtonVariants]: string } = {
  primary: 'bg-blue-600 text-white shadow-sm hover:bg-blue-500  focus-visible:outline-2  focus-visible:outline-blue-600',
  secondary: 'bg-white ring-1 ring-inset ring-gray-300 hover:bg-gray-50',
  danger: 'bg-red-600 text-white focus-visible:outline-red-600',
  warning: 'bg-orange-500 text-white focus-visible:outline-orange-500',
  success: 'bg-green-400 text-white focus-visible:outline-green-400',
};

export const Button = ({ children, HeroIcon, variant, type, className, loading, ...props }: Props) => {
  return (
    <button
      type={type}
      className={classNames(
        'rounded px-2.5 py-1.5 text-sm font-semibold shadow-sm focus-visible:outline focus-visible:outline-offset-2',
        HeroIcon || loading ? 'flex items-center justify-center' : '',
        variants[variant ?? 'primary'],
        className ?? ''
      )}
      {...props}
    >
      {HeroIcon && <HeroIcon className={classNames('-ml-0.5 h-5 w-5', children ? 'mr-0.5' : '')} aria-hidden="true" />}
      {loading && (
        <div
          className="inline-block h-4 w-4 -ml-0.5 mr-2 animate-spin rounded-full border-4 border-solid border-current border-r-transparent align-[-0.125em] text-primary motion-reduce:animate-[spin_1.5s_linear_infinite]"
          role="status"
        >
          <span className="!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]">Loading...</span>
        </div>
      )}
      {children}
    </button>
  );
};
