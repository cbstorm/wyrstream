export default function generate(plop) {
  plop.setGenerator('route', {
    description: 'Generate Route',
    prompts: [
      {
        type: 'input',
        name: 'name',
        message: 'What is the name of dto (singular)?',
      },
    ],
    actions: function (data) {
      const actions = [
        {
          type: 'add',
          path: 'control_service/http_server/{{name}}.route.go',
          templateFile: '__generator/__templates/route.hbs',
        },
      ];
      return actions;
    },
  });
}
