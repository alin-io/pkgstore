import { DiNpm, DiPython, DiDocker } from 'react-icons/di';
import { IconType } from 'react-icons/lib';

interface Props {
  name: string;
  className?: string;
}

const icons: Record<string, IconType> = {
  npm: DiNpm,
  pypi: DiPython,
  container: DiDocker,
};

export function ServiceIcon({ name, className }: Props) {
  const Icon = icons[name] ?? null;
  return Icon ? <Icon className={className} width={30} height={30} /> : <></>;
}
