export default function generate(plop) {
    plop.setGenerator("repository", {
      description: "Generate Entity",
      prompts: [
        {
          type: "input",
          name: "name",
          message: "What is the name of repository (singular)?",
        },
      ],
      actions: function (data) {
        const actions = [
          {
            type: "add",
            path: "lib/repositories/{{name}}.repository.go",
            templateFile: "__generator/__templates/repository.hbs",
          },
        ];
        return actions;
      },
    });
  }