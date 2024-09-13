import dto from './__generator/dto.generator.mjs';
import entity from './__generator/entity.generator.mjs';
import job from './__generator/job.generator.mjs';
import module from './__generator/module.generator.mjs';
import repository from './__generator/repository.generator.mjs';
import route from './__generator/route.generator.mjs';
import service from './__generator/service.generator.mjs';
const generate = (plop) => {
  entity(plop);
  service(plop);
  repository(plop);
  dto(plop);
  module(plop);
  route(plop);
  job(plop);

  plop.setHelper('Case', function (text) {
    return (text[0].toUpperCase() + text.slice(1)).replaceAll('_', '');
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
    if (['y'].includes(text.slice(-1))) {
      return text[0].toUpperCase() + text.slice(1, text.length - 1) + 'ies';
    }
    if (['ch', 'sh', 'ss'].includes(text.slice(-2))) {
      return text[0].toUpperCase() + text.slice(1) + 'es';
    }
    if (['x', 'z', 's'].includes(text.slice(-1))) {
      return text[0].toUpperCase() + text.slice(1) + 'es';
    }
    return text[0].toUpperCase() + text.slice(1) + 's';
  });
  plop.setHelper('CAPP', function (text) {
    return text.toUpperCase();
  });
};

export default generate;
