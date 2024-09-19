import dto from './__generator/dto.generator.mjs';
import entity from './__generator/entity.generator.mjs';
import module from './__generator/module.generator.mjs';
import repository from './__generator/repository.generator.mjs';
import route from './__generator/route.generator.mjs';
import service from './__generator/service.generator.mjs';

const handleCase = (text) => {
  const words = text.split('_');
  for (let i = 0; i < words.length; i++) {
    words[i] = words[i][0].toUpperCase() + words[i].slice(1);
  }
  return words.join('');
};
const generate = (plop) => {
  entity(plop);
  service(plop);
  repository(plop);
  dto(plop);
  module(plop);
  route(plop);

  plop.setHelper('Case', function (text) {
    return handleCase(text);
  });
  plop.setHelper('Plural', function (text) {
    if (['y'].includes(text.slice(-1))) {
      return text.slice(0, text.length - 1) + 'ies';
    }
    if (['ch', 'sh', 'ss'].includes(text.slice(-2))) {
      return text + 'es';
    }
    if (['x', 'z', 's'].includes(text.slice(-1))) {
      return text + 'es';
    }
    return text + 's';
  });
  plop.setHelper('PluralCase', function (text) {
    text = handleCase(text);
    if (['y'].includes(text.slice(-1))) {
      return text[0].toUpperCase() + text.slice(1, text.length - 1) + 'ies';
    }
    if (['ch', 'sh', 'ss'].includes(text.slice(-2))) {
      return text[0].toUpperCase() + text.slice(1) + 'es';
    }
    if (['x', 'z', 's'].includes(text.slice(-1))) {
      return text[0].toUpperCase() + text.slice(1) + 'es';
    }
    return text + 's';
  });
  plop.setHelper('CAPP', function (text) {
    return text.toUpperCase();
  });
};

export default generate;
