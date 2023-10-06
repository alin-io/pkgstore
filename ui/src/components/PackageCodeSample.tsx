import { SERVER_HOST } from '../api';

interface Props {
  name: string;
  service: string;
  type: 'pull' | 'push';
}

export function PackageCodeSample({ type, service, name }: Props) {
  return (
    <div className="w-full prose">
      <pre>{ActionCodeSample(name)?.[type]?.[service] ?? ''}</pre>
    </div>
  );
}

function ActionCodeSample(name: string): Record<string, Record<string, string>> {
  const serverDomain = SERVER_HOST.replace('http://', '').replace('https://', '');
  return {
    pull: {
      container: `docker pull ${serverDomain}/${name}`,
      npm: `npm install ${serverDomain}/${name} --registry=${SERVER_HOST}/npm`,
      pypi: `pip install ${serverDomain}/${name} --index-url=${SERVER_HOST}/pypi`,
    },
    push: {
      container: `docker push ${serverDomain}/${name}`,
      npm: `npm publish --registry=${SERVER_HOST}/npm`,
      pypi: `poetry publish --build --repository=${SERVER_HOST}/pypi`,
    },
  };
}
