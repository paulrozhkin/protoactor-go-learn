# protoactor-go-learn

Тестовый проект для обучения protoactor-go.
В рамках проекта создан:
1. Bank, который порождает дочерних акторов и обеспечивает работу с Bank Account по GUID
2. Взаимодействие с Bank через actor model
3. Поддержка Supervision для инварианта Bank Account (если счет меньше нуля)
4. Внедрен Behavior для осуществления Internal Transaction между двумя Bank Account