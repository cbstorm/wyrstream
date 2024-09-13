export default function generate(plop) {
    plop.setGenerator('entity', {
      description: 'Generate Entity',
      prompts: [
        {
          type: 'input',
          name: 'name',
          message: 'What is the name of entity (singular)?',
        },
      ],
      actions: function (data) {
        const actions = [
          {
            type: 'add',
            path: 'lib/entities/{{name}}.entity.go',
            templateFile: '__generator/__templates/entity.hbs',
          },
        ];
        return actions;
      },
    });
  }